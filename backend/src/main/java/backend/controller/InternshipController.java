package backend.controller;

import backend.dto.AllTechDTO;
import backend.dto.InternshipDTO;
import backend.service.InternshipService;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Set;

@CrossOrigin(origins = {
        "http://localhost:5173",
        "http://localhost:3000"
})@RestController
@RequestMapping("/api/internship/")
@RequiredArgsConstructor
public class InternshipController {

    private final InternshipService internshipService;

    @GetMapping("/all")
    public Page<InternshipDTO> getAllInternships(
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size
    ) {
        return internshipService.getInternships(page, size);
    }

    @GetMapping
    public Page<InternshipDTO> getAllBySpecifications(
            @RequestParam(required = false) Set<String> tech,
            @RequestParam(required = false) String location,
            @RequestParam(required = false) Integer minSalary,
            @RequestParam(required = false) String companyName,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size
    ) {

        return internshipService.searchInternships(tech, minSalary, location, companyName, page, size);

    }

    @GetMapping("/tech")
    public List<AllTechDTO> getTechInternships() {
        return internshipService.getAllTech();
    }

}
