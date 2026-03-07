package backend.service;

import com.aggregator.internship.CompanyInternship;
import org.springframework.stereotype.Service;

@Service
public class InternshipService {

    public void saveFromProto(CompanyInternship internship) {
        System.out.println("Company = " + internship.getCompanyName());
        System.out.println("Source = " + internship.getSourceSite());

        internship.getTracksList().forEach(track -> {
            System.out.println("Position = " + track.getPositionName());
            System.out.println("Tech = " + track.getTechStackList());
        });
    }
}
